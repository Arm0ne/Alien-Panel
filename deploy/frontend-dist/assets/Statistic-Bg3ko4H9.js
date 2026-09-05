import{$n as e,An as t,Ii as n,Jn as r,Nn as i,Ri as a,Xr as o,Yr as s,ki as c,qn as l}from"./router-CLsq84Q3.js";import{n as u}from"./index-ChGrRyuT.js";var d=s(`statistic`,[o(`label`,`
 font-weight: var(--n-label-font-weight);
 transition: .3s color var(--n-bezier);
 font-size: var(--n-label-font-size);
 color: var(--n-label-text-color);
 `),s(`statistic-value`,`
 margin-top: 4px;
 font-weight: var(--n-value-font-weight);
 `,[o(`prefix`,`
 margin: 0 4px 0 0;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-prefix-text-color);
 `,[s(`icon`,{verticalAlign:`-0.125em`})]),o(`content`,`
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-text-color);
 `),o(`suffix`,`
 margin: 0 0 0 4px;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-suffix-text-color);
 `,[s(`icon`,{verticalAlign:`-0.125em`})])])]),f=n({name:`Statistic`,props:Object.assign(Object.assign({},t.props),{tabularNums:Boolean,label:String,value:[String,Number]}),slots:Object,setup(e){let{mergedClsPrefixRef:n,inlineThemeDisabled:a,mergedRtlRef:o}=r(e),s=t(`Statistic`,`-statistic`,d,u,e,n),f=i(`Statistic`,o,n),p=c(()=>{let{self:{labelFontWeight:e,valueFontSize:t,valueFontWeight:n,valuePrefixTextColor:r,labelTextColor:i,valueSuffixTextColor:a,valueTextColor:o,labelFontSize:c},common:{cubicBezierEaseInOut:l}}=s.value;return{"--n-bezier":l,"--n-label-font-size":c,"--n-label-font-weight":e,"--n-label-text-color":i,"--n-value-font-weight":n,"--n-value-font-size":t,"--n-value-prefix-text-color":r,"--n-value-suffix-text-color":a,"--n-value-text-color":o}}),m=a?l(`statistic`,void 0,p,e):void 0;return{rtlEnabled:f,mergedClsPrefix:n,cssVars:a?void 0:p,themeClass:m?.themeClass,onRender:m?.onRender}},render(){var t;let{mergedClsPrefix:n,$slots:{default:r,label:i,prefix:o,suffix:s}}=this;return(t=this.onRender)==null||t.call(this),a(`div`,{class:[`${n}-statistic`,this.themeClass,this.rtlEnabled&&`${n}-statistic--rtl`],style:this.cssVars},e(i,e=>a(`div`,{class:`${n}-statistic__label`},this.label||e)),a(`div`,{class:`${n}-statistic-value`,style:{fontVariantNumeric:this.tabularNums?`tabular-nums`:``}},e(o,e=>e&&a(`span`,{class:`${n}-statistic-value__prefix`},e)),this.value===void 0?e(r,e=>e&&a(`span`,{class:`${n}-statistic-value__content`},e)):a(`span`,{class:`${n}-statistic-value__content`},this.value),e(s,e=>e&&a(`span`,{class:`${n}-statistic-value__suffix`},e))))}});export{f as t};