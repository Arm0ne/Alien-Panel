import{Dr as e,Si as t,Sr as n,Zi as r,aa as i,ar as a,cr as o,oa as s,xi as c,xr as l}from"./router-XnvCoRTR.js";import{n as u}from"./index-CjMER2w_.js";var d=c(`statistic`,[t(`label`,`
 font-weight: var(--n-label-font-weight);
 transition: .3s color var(--n-bezier);
 font-size: var(--n-label-font-size);
 color: var(--n-label-text-color);
 `),c(`statistic-value`,`
 margin-top: 4px;
 font-weight: var(--n-value-font-weight);
 `,[t(`prefix`,`
 margin: 0 4px 0 0;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-prefix-text-color);
 `,[c(`icon`,{verticalAlign:`-0.125em`})]),t(`content`,`
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-text-color);
 `),t(`suffix`,`
 margin: 0 0 0 4px;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-suffix-text-color);
 `,[c(`icon`,{verticalAlign:`-0.125em`})])])]),f=i({name:`Statistic`,props:Object.assign(Object.assign({},a.props),{tabularNums:Boolean,label:String,value:[String,Number]}),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:i,mergedRtlRef:s}=n(e),c=a(`Statistic`,`-statistic`,d,u,e,t),f=o(`Statistic`,s,t),p=r(()=>{let{self:{labelFontWeight:e,valueFontSize:t,valueFontWeight:n,valuePrefixTextColor:r,labelTextColor:i,valueSuffixTextColor:a,valueTextColor:o,labelFontSize:s},common:{cubicBezierEaseInOut:l}}=c.value;return{"--n-bezier":l,"--n-label-font-size":s,"--n-label-font-weight":e,"--n-label-text-color":i,"--n-value-font-weight":n,"--n-value-font-size":t,"--n-value-prefix-text-color":r,"--n-value-suffix-text-color":a,"--n-value-text-color":o}}),m=i?l(`statistic`,void 0,p,e):void 0;return{rtlEnabled:f,mergedClsPrefix:t,cssVars:i?void 0:p,themeClass:m?.themeClass,onRender:m?.onRender}},render(){var t;let{mergedClsPrefix:n,$slots:{default:r,label:i,prefix:a,suffix:o}}=this;return(t=this.onRender)==null||t.call(this),s(`div`,{class:[`${n}-statistic`,this.themeClass,this.rtlEnabled&&`${n}-statistic--rtl`],style:this.cssVars},e(i,e=>s(`div`,{class:`${n}-statistic__label`},this.label||e)),s(`div`,{class:`${n}-statistic-value`,style:{fontVariantNumeric:this.tabularNums?`tabular-nums`:``}},e(a,e=>e&&s(`span`,{class:`${n}-statistic-value__prefix`},e)),this.value===void 0?e(r,e=>e&&s(`span`,{class:`${n}-statistic-value__content`},e)):s(`span`,{class:`${n}-statistic-value__content`},this.value),e(o,e=>e&&s(`span`,{class:`${n}-statistic-value__suffix`},e))))}});export{f as t};