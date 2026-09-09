import{$i as e,Ki as t,Qn as n,br as r,ea as i,hi as a,hr as o,mi as s,mr as c,tr as l}from"./router-C2u90Eon.js";import{n as u}from"./index-CRoWbJ3Q.js";var d=s(`statistic`,[a(`label`,`
 font-weight: var(--n-label-font-weight);
 transition: .3s color var(--n-bezier);
 font-size: var(--n-label-font-size);
 color: var(--n-label-text-color);
 `),s(`statistic-value`,`
 margin-top: 4px;
 font-weight: var(--n-value-font-weight);
 `,[a(`prefix`,`
 margin: 0 4px 0 0;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-prefix-text-color);
 `,[s(`icon`,{verticalAlign:`-0.125em`})]),a(`content`,`
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-text-color);
 `),a(`suffix`,`
 margin: 0 0 0 4px;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-suffix-text-color);
 `,[s(`icon`,{verticalAlign:`-0.125em`})])])]),f=e({name:`Statistic`,props:Object.assign(Object.assign({},n.props),{tabularNums:Boolean,label:String,value:[String,Number]}),slots:Object,setup(e){let{mergedClsPrefixRef:r,inlineThemeDisabled:i,mergedRtlRef:a}=o(e),s=n(`Statistic`,`-statistic`,d,u,e,r),f=l(`Statistic`,a,r),p=t(()=>{let{self:{labelFontWeight:e,valueFontSize:t,valueFontWeight:n,valuePrefixTextColor:r,labelTextColor:i,valueSuffixTextColor:a,valueTextColor:o,labelFontSize:c},common:{cubicBezierEaseInOut:l}}=s.value;return{"--n-bezier":l,"--n-label-font-size":c,"--n-label-font-weight":e,"--n-label-text-color":i,"--n-value-font-weight":n,"--n-value-font-size":t,"--n-value-prefix-text-color":r,"--n-value-suffix-text-color":a,"--n-value-text-color":o}}),m=i?c(`statistic`,void 0,p,e):void 0;return{rtlEnabled:f,mergedClsPrefix:r,cssVars:i?void 0:p,themeClass:m?.themeClass,onRender:m?.onRender}},render(){var e;let{mergedClsPrefix:t,$slots:{default:n,label:a,prefix:o,suffix:s}}=this;return(e=this.onRender)==null||e.call(this),i(`div`,{class:[`${t}-statistic`,this.themeClass,this.rtlEnabled&&`${t}-statistic--rtl`],style:this.cssVars},r(a,e=>i(`div`,{class:`${t}-statistic__label`},this.label||e)),i(`div`,{class:`${t}-statistic-value`,style:{fontVariantNumeric:this.tabularNums?`tabular-nums`:``}},r(o,e=>e&&i(`span`,{class:`${t}-statistic-value__prefix`},e)),this.value===void 0?r(n,e=>e&&i(`span`,{class:`${t}-statistic-value__content`},e)):i(`span`,{class:`${t}-statistic-value__content`},this.value),r(s,e=>e&&i(`span`,{class:`${t}-statistic-value__suffix`},e))))}});export{f as t};